/**
 * Suivi local de la certification JEB par utilisateur.
 *
 * - `initPending(id)` : appelé à l'inscription d'un candidat -> le quiz devient
 *   obligatoire (le garde routeur renvoie sur /quiz tant que `done` est faux).
 * - `complete(id, ...)` : appelé à la validation du quiz -> débloque + mémorise
 *   le badge et le score.
 *
 * Les comptes existants (sans entrée) ne sont pas bloqués.
 */

const KEY = (id: string) => `certif_${id}`

export interface CertState {
  done: boolean
  badge: boolean
  score: number
}

export function getCert(id: string | undefined | null): CertState | null {
  if (!id) return null
  try {
    const raw = localStorage.getItem(KEY(id))
    return raw ? (JSON.parse(raw) as CertState) : null
  } catch {
    return null
  }
}

function write(id: string, state: CertState) {
  try {
    localStorage.setItem(KEY(id), JSON.stringify(state))
  } catch {
    /* stockage indisponible */
  }
}

/** Marque le quiz comme requis mais pas encore passé. */
export function initPending(id: string) {
  if (!getCert(id)) write(id, { done: false, badge: false, score: 0 })
}

/** Le candidat doit-il encore passer le quiz ? */
export function needsQuiz(id: string | undefined | null): boolean {
  const c = getCert(id)
  return c !== null && !c.done
}

/** Enregistre le résultat du quiz et débloque l'accès. */
export function complete(id: string, badge: boolean, score: number) {
  write(id, { done: true, badge, score })
}
