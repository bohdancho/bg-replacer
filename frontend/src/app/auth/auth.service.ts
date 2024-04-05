import { HttpClient, HttpErrorResponse } from '@angular/common/http'
import { Injectable, inject } from '@angular/core'
import { toSignal } from '@angular/core/rxjs-interop'
import { Subject, catchError, of, sample, startWith, take, tap } from 'rxjs'

type User = {
  id: number
  username: string
}

type RegistrationDTO = {
  username: string
  password: string
}
type LoginDTO = RegistrationDTO

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private http = inject(HttpClient)

  private authChange$ = new Subject<void>()

  user$ = this.http.get<User | null>('api/current-user').pipe(
    sample(this.authChange$.pipe(startWith(123))),
    catchError((err: HttpErrorResponse) => {
      if (err.status !== 401) {
        alert('Unexpected authentification error')
      }
      return of(null)
    }),
  )
  user = toSignal(this.user$)

  constructor() {
    this.fetchUser()
  }

  logout() {
    this.http
      .post('api/logout', null)
      .pipe(take(1))
      .subscribe(() => this.authChange$.next())
  }

  login(payload: LoginDTO) {
    return this.http
      .post<void>('api/login', payload, { responseType: 'json' })
      .pipe(tap(() => this.authChange$.next()))
  }

  register(payload: RegistrationDTO) {
    return this.http.post('api/registration', payload, { responseType: 'json' })
  }

  private fetchUser() {}
}
