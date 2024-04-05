import { Component, computed, input } from '@angular/core'
import { toSignal, toObservable } from '@angular/core/rxjs-interop'
import { debounceTime } from 'rxjs'

@Component({
  selector: 'app-img-view',
  standalone: true,
  imports: [],
  host: {
    class: 'flex-1 self-stretch relative',
  },
  template: `
    @if (processingErrorMsg() && processingEnabled()) {
      <div
        class="text-red-500 text-5xl flex justify-center h-full flex-col text-center"
      >
        <p>Processing errored:</p>
        <p>{{ processingErrorMsg() }}</p>
      </div>
    } @else if (shouldDisplaySpinner()) {
      <div class="w-full h-full flex justify-center items-center">
        loading...
      </div>
    } @else {
      <img alt="" [src]="src()" class="absolute w-full h-full object-contain" />
    }
  `,
})
export class ImgViewComponent {
  processingEnabled = input<boolean>()
  processingErrorMsg = input<null | string>()
  originalSrc = input<string | null>(null)
  processedSrc = input<string | null>()

  src = computed(() => {
    const processingEnabled = this.processingEnabled()
    const processedSrc = this.processedSrc()
    const originalSrc = this.originalSrc()

    if (processingEnabled) {
      return processedSrc ?? originalSrc
    } else {
      return originalSrc
    }
  })

  shouldDisplaySpinner = toSignal(
    toObservable(
      computed(() => {
        const processingEnabled = this.processingEnabled()
        const processedSrc = this.processedSrc()
        const error = this.processingErrorMsg()
        return processingEnabled && !processedSrc && !error
      }),
    ).pipe(debounceTime(200)),
  )
}
