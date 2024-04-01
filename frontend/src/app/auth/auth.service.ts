import { HttpClient, HttpErrorResponse } from '@angular/common/http'
import { Injectable, inject, signal } from '@angular/core'
import { toObservable } from '@angular/core/rxjs-interop'
import { take, tap } from 'rxjs'

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

  private _user = signal<User | null>(null)
  user = this._user.asReadonly()
  user$ = toObservable(this._user) // TODO: doesn't work for AuthGuard

  constructor() {
    this.fetchUser()
  }

  logout() {
    this.http
      .post('api/logout', null)
      .pipe(take(1))
      .subscribe(() => this.fetchUser())
  }

  login(payload: LoginDTO) {
    return this.http
      .post<void>('api/login', payload, { responseType: 'json' })
      .pipe(tap(() => this.fetchUser()))
  }

  register(payload: RegistrationDTO) {
    return this.http.post('api/registration', payload, { responseType: 'json' })
  }

  private fetchUser() {
    this.http
      .get<User | null>('api/current-user')
      .pipe(take(1))
      .subscribe({
        next: (user) => this._user.set(user),
        error: (err: HttpErrorResponse) => {
          if (err.status !== 401) {
            alert('Unexpected authentification error')
          }
          return this._user.set(null)
        },
      })
  }
}
