# This memo summarizes the changes implemented in openlab v2
## optimistic processing
below a certain item value providers start rocessing tasks the moment they receive the request via http. If the payment is not received within the escrow contract beyond an ``òptimism_threshold``` the service stops running and communicates its refusal to continue with the client.
