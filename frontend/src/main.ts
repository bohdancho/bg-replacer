import { bootstrapApplication } from '@angular/platform-browser'
import { AppComponent } from './app/app.component'
import { baseUrlInterceptor, baseUrlToken } from './app/base-url-interceptor'
import { provideHttpClient, withInterceptors } from '@angular/common/http'
import { provideRouter } from '@angular/router'
import { environment } from '../environments/environment'
import { routes } from './app/app.routes'

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    {
      provide: baseUrlToken,
      useValue: environment.apiUrl,
    },
    provideHttpClient(withInterceptors([baseUrlInterceptor])),
  ],
}).catch((err) => console.error(err))
