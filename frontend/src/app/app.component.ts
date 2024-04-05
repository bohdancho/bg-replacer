import { HttpClient, HttpErrorResponse } from '@angular/common/http'
import { Component, computed, inject } from '@angular/core'
import { catchError, map } from 'rxjs/operators'
import { of } from 'rxjs'
import { toSignal } from '@angular/core/rxjs-interop'
import { ImageProcessorComponent } from './img-processor'
import { CommonModule } from '@angular/common'
import { ReactiveFormsModule } from '@angular/forms'
import { RouterLink, RouterOutlet } from '@angular/router'
import { AuthService } from './auth/auth.service'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    CommonModule,
    ReactiveFormsModule,
    ImageProcessorComponent,
    RouterLink,
  ],
  template: `
    <main class="flex items-center px-16 py-6 h-dvh flex-col gap-10 min-h-dvh">
      <h1 class="text-indigo-700 text-center text-5xl">
        {{ greeting() }} Welcome to
        <a routerLink="/" class="underline font-bold">imaginaer</a>
      </h1>
      <router-outlet />
      <div class="mt-auto">
        @if (auth.user()) {
          <button (click)="auth.logout()" class="mx-auto block">logout</button>
        } @else {
          <a class="mx-auto block" routerLink="/login">sign in</a>
          <a class="mx-auto block" routerLink="/registration">sign up</a>
        }
        <p>health: {{ health() }}</p>
      </div>
    </main>
  `,
})
export class AppComponent {
  readonly http = inject(HttpClient)
  readonly auth = inject(AuthService)
  health = toSignal(
    this.http.get('api/health').pipe(
      map(() => 'all good'),
      catchError((err: HttpErrorResponse) => {
        return of(`not good, status: ${err.status}`)
      }),
    ),
  )

  greeting = computed(() => {
    const username = this.auth.user()?.username
    return username ? `Hi ${username}!` : 'Hi!'
  })
}
