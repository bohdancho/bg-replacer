import { CommonModule } from '@angular/common'
import { HttpClient, HttpClientModule } from '@angular/common/http'
import { Component, ElementRef, ViewChild, signal, viewChild } from '@angular/core'
import { RouterOutlet } from '@angular/router'
import { from, fromEvent, of } from 'rxjs'
import { filter, map, switchMap, tap } from 'rxjs/operators'
import { toObservable } from '@angular/core/rxjs-interop'

// const PNG_HEADER = 'data:image/png;base64,'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, HttpClientModule, CommonModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'bg-replacer'

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

  constructor(private http: HttpClient) {}
}
