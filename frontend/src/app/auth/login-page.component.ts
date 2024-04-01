import { NgIf } from '@angular/common'
import { HttpErrorResponse } from '@angular/common/http'
import { Component, computed, inject, signal } from '@angular/core'
import { Validators, ReactiveFormsModule, FormBuilder } from '@angular/forms'
import { RouterLink } from '@angular/router'
import { take } from 'rxjs'
import { AuthService } from './auth.service'

@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [ReactiveFormsModule, NgIf, RouterLink],
  template: `
    @if (auth.user()) {
      <div class="flex flex-col gap-4 w-80 max-w-full">
        <h2 class="text-2xl my-0 self-center">Success!</h2>
        <p class="text-center">You are now authorized</p>
      </div>
    } @else {
      <form
        class="flex flex-col gap-4 w-80 max-w-full"
        [formGroup]="form"
        (submit)="onSubmit($event)"
      >
        <h2 class="text-2xl my-0 self-center">Login</h2>
        <input type="text" formControlName="username" placeholder="Username" />
        @if (form.dirty) {
          <div *ngIf="form.controls.username.errors?.['required']">
            Username is required.
          </div>
        }
        <input
          type="password"
          formControlName="password"
          placeholder="Password"
        />
        @if (form.dirty) {
          <div *ngIf="form.controls.password.errors?.['required']">
            Password is required.
          </div>
        }
        <button
          type="submit"
          [disabled]="pending()"
          class="bg-blue-500 text-white border-none rounded-md cursor-pointer disabled:cursor-auto disabled:bg-gray-300 py-2"
        >
          Sign in
        </button>
        <div *ngIf="invalidCredentials()">
          The provided username or password is incorrect.
        </div>
        <div *ngIf="internalServerError()">
          An unexpected server error occured.
        </div>
      </form>
    }
  `,
})
export class LoginPageComponent {
  auth = inject(AuthService)
  fb = inject(FormBuilder)

  pending = signal(false)
  form = this.fb.group(
    {
      username: ['', [Validators.required]],
      password: ['', [Validators.required]],
    },
    { updateOn: 'submit' },
  )

  private responseErrorCode = signal<number | null>(null)
  invalidCredentials = computed(() =>
    this.responseErrorCode() === 401 ? true : false,
  )
  internalServerError = computed(() =>
    this.responseErrorCode() === 500 ? true : false,
  )

  onSubmit(e: Event) {
    e.preventDefault()

    const { username, password } = this.form.value
    if (this.form.valid && username && password) {
      this.pending.set(true)
      this.auth
        .login({ username, password })
        .pipe(take(1))
        .subscribe({
          next: () => {
            this.pending.set(false)
            this.form.reset()
          },
          error: (error: HttpErrorResponse) => {
            this.pending.set(false)
            this.responseErrorCode.set(error.status)
            this.form.valueChanges
              .pipe(take(1))
              .subscribe(() => this.responseErrorCode.set(null))
          },
        })
    }
  }
}
