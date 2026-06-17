export function cronExpr(time: string): string {
  const [hh, mm] = time.split(':')
  const hour = Number(hh)
  const minute = Number(mm)
  return `${minute} ${hour} * * *`
}

export function weeklyCronExpr(day: string, time: string): string {
  const days: Record<string, number> = {
    sunday: 0,
    monday: 1,
    tuesday: 2,
    wednesday: 3,
    thursday: 4,
    friday: 5,
    saturday: 6,
  }
  const dow = days[day.toLowerCase()] ?? 0
  const [hh, mm] = time.split(':')
  return `${Number(mm)} ${Number(hh)} * * ${dow}`
}
