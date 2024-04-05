import { Routes } from '@angular/router'
import { LoginPageComponent } from './auth/login-page.component'
import { ImageProcessorComponent } from './img-processor'
import { HomePageComponent } from './home-page.component'
import { RegistrationPageComponent } from './auth/registration-page.component'
import { isAuthenticatedGuard } from './auth/auth-guard.service'
import { GalleryComponent } from './gallery/gallery.component'

export const routes: Routes = [
  { path: '', component: HomePageComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'registration', component: RegistrationPageComponent },
  { path: 'grayscale', component: ImageProcessorComponent },
  {
    path: 'gallery',
    component: GalleryComponent,
    canActivate: [isAuthenticatedGuard],
  },
  { path: '**', redirectTo: '/' },
]
