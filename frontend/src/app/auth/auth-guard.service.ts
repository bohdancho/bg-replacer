import { inject } from '@angular/core'
import { CanActivateFn, Router } from '@angular/router'
import { AuthService } from './auth.service'
import { map, tap } from 'rxjs'

export const canActivateAuthed: CanActivateFn = () => {
  const auth = inject(AuthService)
  const router = inject(Router)

  return auth.user$.pipe(
    map(Boolean),
    map((isAuthed) => isAuthed || router.createUrlTree(['login'])),
    tap(console.log),
  )
}
