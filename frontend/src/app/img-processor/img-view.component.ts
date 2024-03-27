import { Component, computed, effect, input } from '@angular/core'
import { toSignal, toObservable } from '@angular/core/rxjs-interop'
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'
import { debounceTime } from 'rxjs'

@Component({
  selector: 'app-img-view',
  standalone: true,
  imports: [MatProgressSpinnerModule],
  host: {
    class: 'flex-1 self-stretch relative',
  },
  template: `
    @if (processingErrorMsg()) {
      <div class="text-red-500 text-5xl flex justify-center h-full flex-col text-center">
        <p>Processing errored:</p>
        <p>{{ processingErrorMsg() }}</p>
      </div>
    } @else if (shouldDisplaySpinner()) {
      <div class="w-full h-full flex justify-center items-center"><mat-spinner diameter="50" /></div>
    } @else {
      <img [src]="src()" class="absolute w-full h-full object-contain" />
    }
  `,
})
export class ImgViewComponent {
  shouldDisplayProcessed = input<boolean>()
  processingErrorMsg = input<null | string>()
  originalSrc = input<string | null>(null)
  processedSrc = input<string | null>()

  src = computed(() => {
    const shouldDisplayProcessed = this.shouldDisplayProcessed()
    const processedSrc = this.processedSrc()
    const originalSrc = this.originalSrc()

    if (shouldDisplayProcessed) {
      return processedSrc ?? originalSrc
    } else {
      return originalSrc
    }
  })

  shouldDisplaySpinner = toSignal(
    toObservable(
      computed(() => {
        const shouldProcess = this.shouldDisplayProcessed()
        const processedImgSrc = this.processedSrc()
        return shouldProcess && !processedImgSrc
      }),
    ).pipe(debounceTime(200)),
  )
}
