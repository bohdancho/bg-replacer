import { HttpClient, HttpClientModule, HttpErrorResponse } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { toSignal } from '@angular/core/rxjs-interop'
import { FormControl, ReactiveFormsModule } from '@angular/forms'
import { RouterOutlet } from '@angular/router'
import { filter, take, map, switchMap, catchError, of, Subject } from 'rxjs'
import { MatSlideToggleModule } from '@angular/material/slide-toggle'
import { ImgViewComponent } from './img-view.component'
import { CommonModule } from '@angular/common'

@Component({
  selector: 'app-img-processor',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, CommonModule, MatSlideToggleModule, ReactiveFormsModule, ImgViewComponent],
  host: {
    class: 'flex-1 flex flex-col justify-center items-center self-stretch gap-4',
  },
  template: `
    @if (originalImgSrc()) {
      <mat-slide-toggle [formControl]="processingControl">Grayscale</mat-slide-toggle>
      <app-img-view
        [shouldDisplayProcessed]="!!processing()"
        [originalSrc]="originalImgSrc()"
        [processedSrc]="processedImgSrc() ?? null"
        [processingErrorMsg]="processingErrorMsg()"
      />
    } @else {
      <input type="file" (change)="onOriginalImgChange($event)" />
    }
  `,
})
export class ImageProcessorComponent {
  readonly http = inject(HttpClient)

  onOriginalImgChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (!files || files.length === 0) return

    const reader = new FileReader()
    reader.readAsDataURL(files[0])
    reader.onload = () => {
      this.originalImgSrc.set(reader.result as string)
    }
  }
  originalImgSrc = signal<string | null>(null)

  processingControl = new FormControl(false)
  processing = toSignal(this.processingControl.valueChanges)
  processingRequested$ = this.processingControl.valueChanges.pipe(filter(Boolean), take(1))

  processingErrorMsg = signal<null | string>(null)

  processedImgSrc$ = this.processingRequested$.pipe(
    map(() => this.originalImgSrc()),
    filter(Boolean),
    switchMap((img) => this.http.post<string>('api/grayscale', { img })),
    catchError((err: HttpErrorResponse) => {
      this.processingErrorMsg.set(err.error || err.message)
      return of(null)
    }),
    map((url) => url),
  )
  processedImgSrc = toSignal(this.processedImgSrc$)
}
