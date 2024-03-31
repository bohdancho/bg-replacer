import { Component } from '@angular/core'
import { RouterLink } from '@angular/router'

@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [RouterLink],
  template: `<a routerLink="/grayscale" class="text-blue-800 text-2xl">grayscale</a>`,
})
export class HomePageComponent {}
