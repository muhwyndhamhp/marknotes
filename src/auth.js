import Clerk from '@clerk/clerk-js';

const clerkPubKey = process.env.CLERK_PUBLIC_KEY

window.Clerk = new Clerk(clerkPubKey)

window.Clerk.load().then(() => {
  document.body.addEventListener('htmx:afterRequest', function (event) {
    if (event.detail.xhr.status === 401) {
      event.detail.xhr
      localStorage.setItem('failed-hx-req', event.detail.xhr.responseURL)

      if (window.Clerk) {
        window.Clerk.handleUnauthenticated().then(window.navigateFailedReq())
      } else {
        window.location.reload()
      }
    }
  })
})

window.navigateFailedReq = function() {
      const failedReq = localStorage.getItem('failed-hx-req')
      localStorage.removeItem('failed-hx-req')
      window.location.href = failedReq
}
