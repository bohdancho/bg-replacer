import { Component, Signal, inject, signal } from '@angular/core'
import { FormControl, ReactiveFormsModule } from '@angular/forms'
import { RouterOutlet } from '@angular/router'
import { filter, map, switchMap, distinctUntilChanged } from 'rxjs'
import { MatSlideToggleModule } from '@angular/material/slide-toggle'
import { ImgViewComponent } from './img-view.component'
import { CommonModule } from '@angular/common'
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop'
import { ImgApiService } from '../img-api.service'

type ImgProcessorState = {
  processingEnabled: Signal<boolean | undefined>
  originalSrc: string | null
  originalBlob: Blob | null
  processedSrc: string | null
  processingError: string | null
}

@Component({
  selector: 'app-img-processor',
  standalone: true,
  imports: [RouterOutlet, CommonModule, MatSlideToggleModule, ReactiveFormsModule, ImgViewComponent],
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
  imgApi = inject(ImgApiService)

  processingControl = new FormControl(false)
  state = signal<ImgProcessorState>({
    processingEnabled: toSignal(this.processingControl.valueChanges.pipe(map(Boolean))),
    originalBlob: null,
    originalSrc: null,
    processedSrc: null,
    processingError: null,
  })

  reset() {
    this.processingControl.setValue(false)
    this.state.update((state) => ({ ...state, originalSrc: null, processedSrc: null, processingError: null }))
  }

  onOriginalChange(event: Event) {
    const file = (event.target as HTMLInputElement).files?.item(0)
    if (!file) return

    this.state.update((state) => ({
      ...state,
      originalSrc: URL.createObjectURL(file),
      originalBlob: file,
    }))
  }

  processingRequest$ = this.processingControl.valueChanges.pipe(
    map(() => this.state().originalBlob),
    filter(Boolean),
    distinctUntilChanged(),
    switchMap((blob) => this.imgApi.grayscale(blob)),
  )

  constructor() {
    this.processingRequest$.pipe(takeUntilDestroyed()).subscribe(({ blob, error }) => {
      if (error !== null) {
        return this.state.update((state) => ({ ...state, processingError: error }))
      }
      return this.state.update((state) => ({ ...state, processedSrc: URL.createObjectURL(blob) }))
    })
  }
}
