package templaterender

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	contentrender "github.com/woragis/management/backend/server/internal/content/templaterender"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	"github.com/woragis/management/backend/server/internal/models"
)

type Engine struct {
	content     *contentsvc.Service
	devProjects *devprojectsvc.Service
}

func NewEngine(content *contentsvc.Service, devProjects *devprojectsvc.Service) *Engine {
	return &Engine{content: content, devProjects: devProjects}
}

type RenderInput struct {
	Template *models.MessageTemplate
	Job      *models.ScheduledJob
}

type RenderResult struct {
	Body        string
	Data        map[string]string
	Skipped     bool
	SkipReason  string
	ExternalRef string
}

func (e *Engine) Render(ctx context.Context, in RenderInput) (*RenderResult, error) {
	tpl := in.Template
	ds := ParseDataSource(in.Job.DataSource)
	if ds.Program == "" {
		ds.Program = strings.TrimSpace(tpl.ProgramSlug)
	}
	if ds.Program == "" {
		ds.Program = "custom"
	}

	bindings := MergeBindings(ParseBindings(tpl.Bindings), ds.Program, tpl.Body)
	vars, skip, skipReason, externalRef, err := e.resolveVars(ctx, ds, in.Job)
	if err != nil {
		return nil, err
	}
	if skip {
		return &RenderResult{Skipped: true, SkipReason: skipReason, ExternalRef: externalRef}, nil
	}

	resolved := map[string]string{}
	for placeholder, binding := range bindings {
		resolved[placeholder] = resolveBinding(binding, vars)
	}
	body := RenderBody(tpl.Body, resolved)
	return &RenderResult{Body: body, Data: resolved, ExternalRef: externalRef}, nil
}

func resolveBinding(binding string, vars map[string]string) string {
	binding = strings.TrimSpace(binding)
	if binding == "" {
		return ""
	}
	if v, ok := vars[binding]; ok {
		return v
	}
	parts := strings.SplitN(binding, ".", 2)
	if len(parts) == 2 {
		if v, ok := vars[parts[1]]; ok {
			return v
		}
	}
	return vars[binding]
}

func (e *Engine) resolveVars(ctx context.Context, ds DataSource, job *models.ScheduledJob) (map[string]string, bool, string, string, error) {
	program := strings.TrimSpace(strings.ToLower(ds.Program))
	switch program {
	case "leetcode":
		return e.resolveLeetcode(ctx, ds, job)
	case "project":
		return e.resolveProject(ctx, ds)
	default:
		return map[string]string{}, false, "", "", nil
	}
}

func (e *Engine) resolveLeetcode(ctx context.Context, ds DataSource, job *models.ScheduledJob) (map[string]string, bool, string, string, error) {
	if e.content == nil {
		return nil, true, "content service unavailable", "", nil
	}
	dispatchType := strings.TrimSpace(job.ProgramAction)
	if dispatchType == "" {
		dispatchType = "problem"
	}
	if strings.HasPrefix(dispatchType, "leetcode/") {
		dispatchType = strings.TrimPrefix(dispatchType, "leetcode/")
	}
	res, err := e.content.ResolveDispatchVars(ctx, dispatchType, strings.TrimSpace(ds.Date))
	if err != nil {
		return nil, false, "", "", err
	}
	if res.Skip {
		return nil, true, res.SkipReason, res.VideoID, nil
	}
	return FromContentVars(res.Vars), false, "", res.VideoID, nil
}

func FromContentVars(v contentrender.Vars) map[string]string {
	return map[string]string{
		"leetcode.number":           v.Number,
		"leetcode.track":            v.Track,
		"leetcode.problemTitle":     v.ProblemTitle,
		"leetcode.difficulty":       v.Difficulty,
		"leetcode.leetcodeUrl":      v.LeetcodeURL,
		"leetcode.submissionUrl":    v.SubmissionURL,
		"leetcode.shortDescription": v.ShortDescription,
		"leetcode.youtubeUrl":       v.YoutubeURL,
		"leetcode.date":             v.Date,
		"leetcode.problemList":      v.ProblemList,
		"leetcode.nextTheme":        v.NextTheme,
		"leetcode.inviteLink":       v.InviteLink,
		"number":                    v.Number,
		"track":                     v.Track,
		"problemTitle":              v.ProblemTitle,
		"difficulty":                v.Difficulty,
		"leetcodeUrl":               v.LeetcodeURL,
		"submissionUrl":             v.SubmissionURL,
		"shortDescription":          v.ShortDescription,
		"youtubeUrl":                v.YoutubeURL,
		"date":                      v.Date,
		"problemList":               v.ProblemList,
		"nextTheme":                 v.NextTheme,
		"inviteLink":                v.InviteLink,
	}
}

func (e *Engine) resolveProject(ctx context.Context, ds DataSource) (map[string]string, bool, string, string, error) {
	if e.devProjects == nil {
		return nil, true, "projects service unavailable", "", nil
	}
	var projectID uuid.UUID
	var err error
	if strings.TrimSpace(ds.ProjectID) != "" {
		projectID, err = uuid.Parse(strings.TrimSpace(ds.ProjectID))
		if err != nil {
			return nil, true, "invalid projectId", "", nil
		}
	} else if strings.TrimSpace(ds.ProjectSlug) != "" {
		p, err := e.devProjects.GetBySlug(ctx, strings.TrimSpace(ds.ProjectSlug))
		if err != nil {
			return nil, true, "project not found", "", nil
		}
		return projectVars(p), false, "", p.ID.String(), nil
	} else {
		return nil, true, "projectId or projectSlug required in dataSource", "", nil
	}
	p, err := e.devProjects.GetByID(ctx, projectID)
	if err != nil {
		return nil, true, "project not found", "", nil
	}
	return projectVars(p), false, "", projectID.String(), nil
}

func projectVars(p *models.Project) map[string]string {
	stack := parseStackJSON(p.Stack)
	return map[string]string{
		"project.name":             p.Name,
		"project.slug":             p.Slug,
		"project.shortDescription": p.ShortDescription,
		"project.demoUrl":          p.DemoURL,
		"project.githubUrl":        p.GithubURL,
		"project.repoUrl":          p.RepoURL,
		"project.stack":            strings.Join(stack, ", "),
		"projectName":              p.Name,
		"projectSlug":              p.Slug,
		"shortDescription":         p.ShortDescription,
		"demoUrl":                  p.DemoURL,
		"githubUrl":                p.GithubURL,
		"repoUrl":                  p.RepoURL,
		"stack":                    strings.Join(stack, ", "),
	}
}

func parseStackJSON(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var out []string
	_ = json.Unmarshal(raw, &out)
	return out
}
