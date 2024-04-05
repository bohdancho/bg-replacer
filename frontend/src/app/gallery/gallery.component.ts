import { HttpClient } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { takeUntilDestroyed } from '@angular/core/rxjs-interop'
import { Subject, merge, startWith, switchMap, tap } from 'rxjs'

@Component({
  selector: 'app-gallery',
  standalone: true,
  imports: [],
  template: `
    @if (state().pending) {
      <p>Loading...</p>
    } @else if (state().src) {
      <img [src]="state().src" alt="gallery image" />
      <button color="primary" (click)="onDelete()">Delete</button>
    } @else {
      <input type="file" (change)="onUpload($event)" />
    }
  `,
})
export class GalleryComponent {
  http = inject(HttpClient)

  state = signal<{
    pending: boolean
    src: string | null
  }>({
    pending: false,
    src: null,
  })

  upload$ = new Subject<Blob>()
  delete$ = new Subject<void>()

  uploaded$ = this.upload$.pipe(
    tap(() => this.state.update((state) => ({ ...state, pending: true }))),
    switchMap((file) => this.http.post<{ url: string }>('api/gallery', file)),
    takeUntilDestroyed(),
  )
  deleted$ = this.delete$.pipe(
    tap(() => this.state.update((state) => ({ ...state, pending: true }))),
    switchMap(() => this.http.delete('api/gallery')),
    takeUntilDestroyed(),
  )

  constructor() {
    merge(this.uploaded$, this.deleted$)
      .pipe(
        startWith(null),
        switchMap(() => this.http.get<{ url: string }>('api/gallery')),
        takeUntilDestroyed(),
      )
      .subscribe(({ url }) =>
        this.state.update(() => ({ pending: false, src: url })),
      )
  }

  onUpload(event: Event) {
    const file = (event.target as HTMLInputElement).files?.item(0)
    if (!file) return

    this.upload$.next(file)
  }

  onDelete() {
    this.delete$.next()
  }
}
