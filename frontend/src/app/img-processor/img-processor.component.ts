import { HttpClient, HttpClientModule, HttpErrorResponse } from '@angular/common/http'
import { Component, Signal, effect, inject, signal } from '@angular/core'
import { FormControl, ReactiveFormsModule } from '@angular/forms'
import { RouterOutlet } from '@angular/router'
import { filter, map, switchMap, catchError, of, distinctUntilChanged, tap } from 'rxjs'
import { MatSlideToggleModule } from '@angular/material/slide-toggle'
import { ImgViewComponent } from './img-view.component'
import { CommonModule } from '@angular/common'
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop'

type ImgProcessorState = {
  processingEnabled: Signal<boolean | undefined>
  originalSrc: string | null
  processedSrc: string | null
  processingError: string | null
}

@Component({
  selector: 'app-img-processor',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, CommonModule, MatSlideToggleModule, ReactiveFormsModule, ImgViewComponent],
  host: {
    class: 'flex-1 flex flex-col justify-center items-center self-stretch gap-4',
  },
  template: `
    @if (state().originalSrc) {
      <mat-slide-toggle [formControl]="processingControl">Grayscale</mat-slide-toggle>
      <app-img-view
        [processingEnabled]="(processingControl.valueChanges | async) ?? false"
        [originalSrc]="state().originalSrc"
        [processedSrc]="state().processedSrc ?? null"
        [processingErrorMsg]="state().processingError ?? null"
      />
      <button mat-raised-button color="primary" (click)="reset()">Reset</button>
    } @else {
      <input type="file" (change)="onOriginalChange($event)" />
    }
  `,
})
export class ImageProcessorComponent {
  readonly http = inject(HttpClient)

  processingControl = new FormControl(false)
  state = signal<ImgProcessorState>({
    processingEnabled: toSignal(this.processingControl.valueChanges.pipe(map(Boolean))),
    originalSrc: null,
    processedSrc: null,
    processingError: null,
  })

  reset() {
    this.processingControl.setValue(false)
    this.state.update((state) => ({ ...state, originalSrc: null, processedSrc: null, processingError: null }))
  }

  onOriginalChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (!files || files.length === 0) return

    const reader = new FileReader()
    reader.readAsDataURL(files[0])
    reader.onload = () => {
      this.state.update((state) => ({
        ...state,
        originalSrc: reader.result as string,
      }))
    }
  }

  processingRequest$ = this.processingControl.valueChanges.pipe(
    tap(console.log),
    filter(Boolean),
    map(() => this.state().originalSrc),
    distinctUntilChanged(),
    switchMap((img) =>
      this.http.post<string>('api/grayscale', { img }).pipe(
        catchError((err) => {
          this.handleProcessingError(err)
          return of(null)
        }),
      ),
    ),
  )

  handleProcessingError(err: HttpErrorResponse) {
    let msg = 'Unknown error'
    if (typeof err.error === 'string') {
      msg = err.error
    } else if (err.message) {
      msg = err.message
    }
    this.state.update((state) => ({ ...state, processingError: msg }))
  }

  constructor() {
    this.processingRequest$.pipe(takeUntilDestroyed()).subscribe((src) => {
      return this.state.update((state) => ({ ...state, processedSrc: src }))
    })

    effect(() => {
      console.log(this.state())
    })
  }
}
