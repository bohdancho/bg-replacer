import { NgIf } from '@angular/common'
import { HttpClient, HttpErrorResponse } from '@angular/common/http'
import { Component, inject, signal } from '@angular/core'
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms'
import { RouterLink } from '@angular/router'
import { take } from 'rxjs'

@Component({
  selector: 'app-registration-page',
  standalone: true,
  imports: [ReactiveFormsModule, NgIf, RouterLink],
  host: { class: 'contents' },
  template: `
    @if (success()) {
      <div class="flex flex-col gap-4 w-80 max-w-full">
        <h2 class="text-2xl my-0 self-center">Success!</h2>
        <p class="text-center">
          You have successfully registered and can now <a routerLink="/login" class="text-blue-800">sign in</a>
        </p>
      </div>
    } @else {
      <form class="flex flex-col gap-4 w-80 max-w-full" [formGroup]="form" (submit)="onSubmit($event)">
        <h2 class="text-2xl my-0 self-center">Sign up</h2>
        <input type="text" formControlName="username" placeholder="Username" />
        @if (form.dirty) {
          <div *ngIf="form.controls.username.errors?.['required']">Username is required.</div>
          <div *ngIf="form.controls.username.errors?.['minlength']">Username must be at least 3 characters long.</div>
          <div *ngIf="form.controls.username.errors?.['maxlength']">Username can be at most 24 characters long.</div>
        }
        <input type="password" formControlName="password" placeholder="Password" />
        @if (form.dirty) {
          <div *ngIf="form.controls.password.errors?.['required']">Password is required.</div>
          <div *ngIf="form.controls.password.errors?.['minlength']">Password must be at least 8 characters long.</div>
          <div *ngIf="form.controls.password.errors?.['maxlength']">Password can be at most 40 characters long.</div>
        }
        <button
          type="submit"
          [disabled]="!this.form.valid"
          class="bg-blue-500 text-white border-none rounded-md cursor-pointer disabled:cursor-auto disabled:bg-gray-300 py-2"
        >
          Sign up
        </button>
        <div *ngIf="errorFromServer">{{ errorFromServer() }}</div>
      </form>
    }
  `,
})
export class RegistrationPageComponent {
  http = inject(HttpClient)

  form = new FormGroup({
    username: new FormControl('', [Validators.required, Validators.minLength(3), Validators.maxLength(24)]),
    password: new FormControl('', [Validators.required, Validators.minLength(8), Validators.maxLength(40)]),
  })

  success = signal(false)
  errorFromServer = signal<string | null>(null)

  onSubmit(e: Event) {
    e.preventDefault()
    if (this.form.valid) {
      this.http.post('api/registration', this.form.value, { responseType: 'json' }).subscribe({
        next: () => this.success.set(true),
        error: (error: HttpErrorResponse) => {
          this.errorFromServer.set(error.statusText)
          this.form.valueChanges.pipe(take(1)).subscribe(() => this.errorFromServer.set(null))
        },
      })
    }
  }
}
