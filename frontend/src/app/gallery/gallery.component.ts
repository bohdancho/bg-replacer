import { HttpClient } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { takeUntilDestroyed } from '@angular/core/rxjs-interop'
import { Subject, catchError, merge, of, startWith, switchMap } from 'rxjs'

const API_PATH = 'api/gallery/'

@Component({
  selector: 'app-gallery',
  standalone: true,
  imports: [],
  template: `
    @if (state().pending) {
      <p>Loading...</p>
    } @else if (state().src) {
      <img [src]="state().src" alt="gallery image" />
      <button color="primary" (click)="delete$.next()">Delete</button>
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
    switchMap((file) => this.http.post<{ url: string }>(API_PATH, file)),
  )
  deleted$ = this.delete$.pipe(switchMap(() => this.http.delete(API_PATH)))

  constructor() {
    merge(this.upload$, this.delete$)
      .pipe(takeUntilDestroyed())
      .subscribe(() => this.state.update(() => ({ pending: true, src: null })))

    merge(this.uploaded$, this.deleted$)
      .pipe(
        startWith(null),
        switchMap(() => this.fetchImage()),
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

  fetchImage() {
    return this.http
      .get<{ url: string }>(API_PATH)
      .pipe(catchError(() => of({ url: null })))
  }
}
