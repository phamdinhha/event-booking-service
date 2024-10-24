# Building a global event booking service
In this repo, I will discuss about how to build a global event booking service. Some of the requirements are:
- Allow users to search and book events worldwide, ensuring a smooth booking experience even under high traffic
- Platform should be able to handle a large number of users, especially during peak hours close to event dates
- Booking process should integrate with a payment gateway (e.g., Stripe) to handle ticket purchases securely
- Upon successful booking, the system should send confirmation and reminder notifications to the users before the event takes place

## Understand the problem
Assume that the system can host 10000 events, each event can have 50000 tickets available, the tickets for each event are sold in one month before the event takes place. Lets say hot event will have 10000 tickets sold in the first 10 minutes.
- User can book tickets from mobile or web app
- User can cancel the reservation

## Non-functional requirements
- High concurrency: During the first 10 minutes for the hot event, the system will have max 2000 requests per second at peak
- Fault tolerance: The system should be able to handle failures gracefully
- Scalability: The system should be able to scale horizontally to handle more traffic
- Moderate latency: The system should be able to handle the requests within a few seconds

## Back of the envelope calculation
- Assume that the hot event will have 50% (25000 tickets) of the total tickets sold in the first 10 minutes. This means that the system will have 2500 successful transactions per minute at peak or 42 successful transactions per second (TPS)
- Assume that the request rate will be higher due to get requests and other requests and drop to around 2000 requests per second at peak.

## Architecture

