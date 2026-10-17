# Feedback reminders

The backend checks reservations every minute and sends a feedback email once the
scheduled reservation date and time are at least six hours in the past.
Cancelled reservations are skipped, and a successful send is recorded so later
checks do not send the same reminder again.

Configure these environment variables:

- `FEEDBACK_PAGE_URL`: absolute URL of the frontend feedback page, for example
  `https://your-frontend.example.com/feedback`. The backend adds the
  `reservation_id` query parameter to the link.
- `RESERVATION_TIMEZONE`: IANA time zone used for reservation date/time values,
  for example `Asia/Kolkata`. If omitted, the server's local time zone is used.

Reservation dates must use `YYYY-MM-DD`. Reservation times may use 24-hour
`HH:MM` or `HH:MM:SS`, or 12-hour `h:MM AM/PM` formats.
