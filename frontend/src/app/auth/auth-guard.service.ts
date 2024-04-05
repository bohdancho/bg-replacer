import { inject } from '@angular/core'
import { CanActivateFn, Router } from '@angular/router'
import { AuthService } from './auth.service'
import { map } from 'rxjs'
import { toObservable } from '@angular/core/rxjs-interop'

export const isAuthenticatedGuard: CanActivateFn = () => {
  const auth = inject(AuthService)
  const router = inject(Router)

  return toObservable(auth.user).pipe(
    map((user) => {
      console.log('user', user)
      if (user) return true
      else return router.createUrlTree(['login'])
    }),
  )
}
