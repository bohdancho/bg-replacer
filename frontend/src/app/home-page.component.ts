import { Component } from '@angular/core'
import { RouterLink } from '@angular/router'

@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [RouterLink],
  template: `<div class="flex flex-col text-center gap-2">
    <a routerLink="/grayscale" class="text-blue-800 text-2xl">grayscale</a>
    <a routerLink="/gallery" class="text-blue-800 text-2xl">gallery</a>
  </div>`,
})
export class HomePageComponent {}
