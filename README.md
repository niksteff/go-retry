# Retry 

This is an example implementation of a retry pkg for basically any function 
you want to retry multiple times.

It supports two backoff strategies which are implemented using an interface.

- linear backoff with a start wait time and a max try ceiling
- exponential backoff with a start and max wait time and a max try ceiling

