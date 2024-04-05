import { Routes } from '@angular/router'
import { LoginPageComponent } from './auth/login-page.component'
import { ImageProcessorComponent } from './img-processor'
import { HomePageComponent } from './home-page.component'
import { RegistrationPageComponent } from './auth/registration-page.component'
import { isAuthenticatedGuard } from './auth/auth-guard.service'

export const routes: Routes = [
  { path: '', component: HomePageComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'registration', component: RegistrationPageComponent },
  { path: 'grayscale', component: ImageProcessorComponent },
  {
    path: 'protected',
    component: ImageProcessorComponent,
    canActivate: [isAuthenticatedGuard],
  },
  { path: '**', redirectTo: '/' },
]
