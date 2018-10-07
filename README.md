# url-shortener

Gokit

- Create circle CI build file 1 hour
- Update Readme 0.5 hour

How to structure the project
 - Use monorepo for ease of review

Design Decision
- Use BDD for test

- Tech
    - Upload file logic
- Listen to goverb
- It uses Base62 algorithm to hash the url

- TODO
 - circle ci
Improvement
 - Add caching layer to improve lookup time like Redis

 Improvement:
 - Try other algorithms (base62 for example)

## How to generate
There are a few ways to generate such as base64. The way I chose is

The main drawback might be but it should be good enough in this scenario
