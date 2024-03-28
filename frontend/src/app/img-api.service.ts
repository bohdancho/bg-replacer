import { HttpClient, HttpErrorResponse } from '@angular/common/http'
import { Injectable, inject } from '@angular/core'
import { Observable, catchError, map, of } from 'rxjs'

type ImgResponse = { blob: Blob; error: null } | { blob: null; error: string }

@Injectable({
  providedIn: 'root',
})
export class ImgApiService {
  readonly http = inject(HttpClient)
  grayscale(blob: Blob): Observable<ImgResponse> {
    return this.http.post('api/grayscale', blob, { responseType: 'blob' }).pipe(
      map((blob) => ({ blob, error: null })),
      catchError(({ error, message }: HttpErrorResponse) =>
        of({ blob: null, error: this.blobToString(error) || message }),
      ),
    )
  }

  private blobToString(blob: Blob): string {
    const url = URL.createObjectURL(blob)
    const xmlRequest = new XMLHttpRequest()
    xmlRequest.open('GET', url, false)
    xmlRequest.send()
    URL.revokeObjectURL(url)
    return xmlRequest.responseText
  }
}
