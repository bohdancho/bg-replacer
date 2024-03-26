import { CommonModule } from '@angular/common'
import { HttpClient, HttpClientModule } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { RouterOutlet } from '@angular/router'
import { filter, map, switchMap } from 'rxjs/operators'
import { toObservable } from '@angular/core/rxjs-interop'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, CommonModule],
  template: `
    <main class="flex items-center p-16 h-dvh flex-col gap-10">
      <h1 class="text-indigo-700 text-5xl">Yooooo man! Welcome to <span class="font-bold">bg-replacer</span></h1>
      <input type="file" class="file-upload" (change)="onFileChange($event)" class="shrink-0" />
      <div class="flex-1 w-full relative">
        <img
          [src]="responseImgUrl$ | async"
          #fileInput
          alt="Uploaded image"
          class="absolute w-full h-full object-contain"
        />
      </div>
    </main>
  `,
})
export class AppComponent {
  http = inject(HttpClient)
  uploadedImgUrl = signal<string | null>(null)
  responseImgUrl$ = toObservable(this.uploadedImgUrl).pipe(
    filter(Boolean),
    switchMap((img) => this.http.post<string>('api/grayscale', { img })),
    map((url) => url),
  )

  onFileChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (!files || files.length === 0) return

    const reader = new FileReader()
    reader.readAsDataURL(files[0])
    reader.onload = () => {
      this.uploadedImgUrl.set(reader.result as string)
    }
  }
}
