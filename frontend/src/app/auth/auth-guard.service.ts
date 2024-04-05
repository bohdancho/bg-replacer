import { inject } from '@angular/core'
import { CanActivateFn, Router } from '@angular/router'
import { AuthService } from './auth.service'
import { map } from 'rxjs'

export const isAuthenticatedGuard: CanActivateFn = () => {
  const auth = inject(AuthService)
  const router = inject(Router)

  return auth.user$.pipe(
    map((user) => {
      if (user) return true
      else return router.createUrlTree(['login'])
    }),
  )
}
