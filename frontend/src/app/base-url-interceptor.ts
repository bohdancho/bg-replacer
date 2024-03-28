import { HttpInterceptorFn } from '@angular/common/http'
import { InjectionToken, inject } from '@angular/core'

export const baseUrlToken = new InjectionToken<string>('BASE_API_URL')

export const baseUrlInterceptor: HttpInterceptorFn = (request, next) => {
  const baseUrl = inject(baseUrlToken)
  console.log(request)
  const apiReq = request.clone({ url: `${baseUrl}${request.url}` })
  return next(apiReq)
}
