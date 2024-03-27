import { HttpClient, HttpClientModule, HttpErrorResponse } from '@angular/common/http'
import { Component, inject } from '@angular/core'
import { catchError, map } from 'rxjs/operators'
import { of } from 'rxjs'
import { toSignal } from '@angular/core/rxjs-interop'
import { ImageProcessorComponent } from './img-processor'
import { CommonModule } from '@angular/common'
import { ReactiveFormsModule } from '@angular/forms'
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'
import { MatSlideToggleModule } from '@angular/material/slide-toggle'
import { RouterOutlet } from '@angular/router'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    HttpClientModule,
    CommonModule,
    MatSlideToggleModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    ImageProcessorComponent,
  ],
  template: `
    <main class="flex items-center p-16 h-dvh flex-col gap-10 min-h-dvh">
      <h1 class="text-indigo-700 text-5xl">Yooooo man! Welcome to <span class="font-bold">imaginaer</span></h1>
      <app-img-processor />
      <p class="mt-auto">health: {{ health() }}</p>
    </main>
  `,
})
export class AppComponent {
  readonly http = inject(HttpClient)
  health = toSignal(
    this.http.get('api/health').pipe(
      map(() => 'all good'),
      catchError((err: HttpErrorResponse) => of(`not good, status: ${err.status}`)),
    ),
  )
}
