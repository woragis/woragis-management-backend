function headers(cfg) {
    return {
        'Content-Type': 'application/json',
        'X-Agent-Key': cfg.agentApiKey,
        Authorization: `Bearer ${cfg.agentApiKey}`,
    };
}
async function request(cfg, path, init) {
    const res = await fetch(`${cfg.managementApiUrl}${path}`, {
        ...init,
        headers: { ...headers(cfg), ...(init?.headers ?? {}) },
    });
    const text = await res.text();
    if (!res.ok) {
        throw new Error(`API ${path} failed: ${res.status} ${text}`);
    }
    if (!text)
        return undefined;
    return JSON.parse(text);
}
export class ManagementClient {
    cfg;
    constructor(cfg) {
        this.cfg = cfg;
    }
    getPersonality() {
        return request(this.cfg, '/v1/internal/agent/personality');
    }
    updatePersonality(patch) {
        return request(this.cfg, '/v1/internal/agent/personality', {
            method: 'PATCH',
            body: JSON.stringify(patch),
        });
    }
    resetPersonality() {
        return request(this.cfg, '/v1/internal/agent/personality/reset', {
            method: 'POST',
        });
    }
    searchContacts(params) {
        const q = new URLSearchParams(params).toString();
        return request(this.cfg, `/v1/internal/agent/tools/contacts?${q}`);
    }
    getContact(id) {
        return request(this.cfg, `/v1/internal/agent/tools/contacts/${id}`);
    }
    createContact(body) {
        return request(this.cfg, '/v1/internal/agent/tools/contacts', {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }
    updateContact(id, body) {
        return request(this.cfg, `/v1/internal/agent/tools/contacts/${id}`, {
            method: 'PATCH',
            body: JSON.stringify(body),
        });
    }
    logInteraction(contactId, body) {
        return request(this.cfg, `/v1/internal/agent/tools/contacts/${contactId}/interactions`, {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }
    listContactsDueFollowUp(before) {
        const q = before ? `?before=${encodeURIComponent(before)}` : '';
        return request(this.cfg, `/v1/internal/agent/tools/contacts/due-follow-up${q}`);
    }
    getContactFinance(id) {
        return request(this.cfg, `/v1/internal/agent/tools/contacts/${id}/finance`);
    }
    listProjects(params = {}) {
        const q = new URLSearchParams(params).toString();
        const suffix = q ? `?${q}` : '';
        return request(this.cfg, `/v1/internal/agent/tools/projects${suffix}`);
    }
    getProject(id) {
        return request(this.cfg, `/v1/internal/agent/tools/projects/${id}`);
    }
    createProject(body) {
        return request(this.cfg, '/v1/internal/agent/tools/projects', {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }
    financeDashboard() {
        return request(this.cfg, '/v1/internal/agent/tools/finance/dashboard');
    }
    financeSummary(year, month) {
        const params = new URLSearchParams();
        if (year)
            params.set('year', String(year));
        if (month)
            params.set('month', String(month));
        const q = params.toString();
        return request(this.cfg, `/v1/internal/agent/tools/finance/summary${q ? `?${q}` : ''}`);
    }
    financeCalendar(year, month) {
        const params = new URLSearchParams();
        if (year)
            params.set('year', String(year));
        if (month)
            params.set('month', String(month));
        const q = params.toString();
        return request(this.cfg, `/v1/internal/agent/tools/finance/calendar${q ? `?${q}` : ''}`);
    }
    listIncomeSources(params = {}) {
        const q = new URLSearchParams(params).toString();
        const suffix = q ? `?${q}` : '';
        return request(this.cfg, `/v1/internal/agent/tools/finance/income-sources${suffix}`);
    }
    listTransactions(params = {}) {
        const q = new URLSearchParams(params).toString();
        const suffix = q ? `?${q}` : '';
        return request(this.cfg, `/v1/internal/agent/tools/finance/transactions${suffix}`);
    }
    createTransaction(body) {
        return request(this.cfg, '/v1/internal/agent/tools/finance/transactions', {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }
    createIncomeSource(body) {
        return request(this.cfg, '/v1/internal/agent/tools/finance/income-sources', {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }
}
