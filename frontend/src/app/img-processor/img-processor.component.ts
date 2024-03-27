import { HttpClient, HttpClientModule, HttpErrorResponse } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { FormControl, ReactiveFormsModule } from '@angular/forms'
import { RouterOutlet } from '@angular/router'
import { filter, take, map, switchMap, catchError, of, Subject, ignoreElements } from 'rxjs'
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
        [processingEnabled]="(processingControl.valueChanges | async) ?? false"
        [originalSrc]="originalImgSrc()"
        [processedSrc]="(processedImgSrc | async) ?? null"
        [processingErrorMsg]="(processingErrorMsg | async) ?? null"
      />
    } @else {
      <input type="file" (change)="onOriginalImgChange($event)" />
    }
  `,
})
export class ImageProcessorComponent {
  readonly http = inject(HttpClient)

  reset = new Subject<void>()
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
  processingEnabled = signal(() => this.processingControl.value)
  processingRequest$ = this.processingControl.valueChanges.pipe(
    filter(Boolean),
    take(1),
    map(() => this.originalImgSrc()),
    filter(Boolean),
    switchMap((img) => this.http.post<string>('api/grayscale', { img })),
  )

  processedImgSrc = this.processingRequest$.pipe(catchError(() => of(null)))
  processingErrorMsg = this.processingRequest$.pipe(
    ignoreElements(),
    catchError((err: HttpErrorResponse) => {
      if (typeof err.error === 'string') {
        return of(err.error)
      }
      if (err.message) {
        return of(err.message)
      }
      return of('Unknown error')
    }),
  )
}
