import { CommonModule } from '@angular/common'
import { HttpClient, HttpClientModule, HttpErrorResponse } from '@angular/common/http'
import { Component, computed, inject, signal } from '@angular/core'
import { RouterOutlet } from '@angular/router'
import { catchError, debounceTime, filter, map, switchMap, take } from 'rxjs/operators'
import { of } from 'rxjs'
import { toObservable, toSignal } from '@angular/core/rxjs-interop'
import { MatSlideToggleModule } from '@angular/material/slide-toggle'
import { FormControl, ReactiveFormsModule } from '@angular/forms'
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'

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
  ],
  template: `
    <main class="flex items-center p-16 h-dvh flex-col gap-10 min-h-dvh">
      <h1 class="text-indigo-700 text-5xl">Yooooo man! Welcome to <span class="font-bold">imaginaer</span></h1>
      @if (originalImgUrl()) {
        <mat-slide-toggle [formControl]="shouldDisplayProcessedControl">Grayscale</mat-slide-toggle>
      } @else {
        <div class="flex flex-1 justify-center items-center text-center">
          <input type="file" (change)="onOriginalImgChange($event)" />
        </div>
      }
      @if (isProcessedImgLoading()) {
        <mat-spinner class="flex-1" diameter="50"></mat-spinner>
      } @else if (visibleImgSrc()) {
        <div class="flex-1 w-full relative">
          <img [src]="visibleImgSrc()" class="absolute w-full h-full object-contain" />
        </div>
      }
      <p class="mt-auto">health: {{ health$ | async }}</p>
    </main>
  `,
})
export class AppComponent {
  readonly http = inject(HttpClient)

  onOriginalImgChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (!files || files.length === 0) return

    const reader = new FileReader()
    reader.readAsDataURL(files[0])
    reader.onload = () => {
      this.originalImgUrl.set(reader.result as string)
    }
  }
  originalImgUrl = signal<string | null>(null)

  shouldDisplayProcessedControl = new FormControl(false)
  shouldDisplayProcessed = toSignal(this.shouldDisplayProcessedControl.valueChanges)
  processingRequested$ = this.shouldDisplayProcessedControl.valueChanges.pipe(filter(Boolean), take(1))

  processedImgSrc$ = this.processingRequested$.pipe(
    map(() => this.originalImgUrl()),
    filter(Boolean),
    switchMap((img) => this.http.post<string>('api/grayscale', { img })),
    map((url) => url),
  )
  processedImgSrc = toSignal(this.processedImgSrc$)

  isProcessedImgLoading = toSignal(
    toObservable(
      computed(() => {
        const shouldProcess = this.shouldDisplayProcessed()
        const processedImgSrc = this.processedImgSrc()
        return shouldProcess && !processedImgSrc
      }),
    ).pipe(debounceTime(200)),
  )

  visibleImgSrc = computed(() => {
    const original = this.originalImgUrl()
    const processed = this.processedImgSrc()
    return this.shouldDisplayProcessed() && processed ? processed : original
  })

  health$ = this.http.get('api/health').pipe(
    map(() => 'all good'),
    catchError((err: HttpErrorResponse) => of(`not good, status: ${err.status}`)),
  )
}
