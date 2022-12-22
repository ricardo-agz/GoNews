# GoNews
CIS 193 Final Project

---

## Description
This is a CRUD REST API for a blog-style website where users can create posts and add hashtags.
By adding #somehashtag to the post content, 'somehashtag' will be added to the post tags and a tag 
named 'somehashtag' will be created with an array of post IDs. 

--- 

## ENV file
This project requires a .env file, a sample .env is included below:
```
MONGODB_DATABASE=my-gonews-db-name
MONGODB_USERNAME=your-mongodb-username
MONGODB_PASSWORD=your-mongodb-password
MONGODB_URL=mongodb+srv://....
```

--- 

## Usage
A sample request to create a user is included below:

```POST http://localhost:8000/users```

```
{
    "Username": "myusername",
    "Email": "myusername@email.co",
    "Password": "123456"
}
```


A sample request to create a post is included below:

```GET http://localhost:8000/users/myusername/posts```

```
{
    "Content": "this is my post #mytag #anothertag"
}
```

--- 

## API
#### GET    /                       
* Home page
#### GET    /users                  
* Returns a list of all users
#### GET    /users/:username        
* Returns user with specified username
#### POST   /users                  
* Creates a new user with the data passed in through the JSON body of the request
#### PUT    /users/:username        
* Updates a user with the new data passed in through the JSON body of the request
#### DELETE /users/:username        
* Deletes the user with the specified username
#### GET    /posts                  
* Returns a list of all posts
#### GET    /users/:username/posts  
* Returns all posts belonging to a specific user
#### GET    /posts/:id              
* Returns post with specified ID
#### POST   /users/:username/posts   
* Creates a new post belonging to user with given username
#### DELETE /users/:username/posts/:id
* Deletes the post with the specified ID
#### GET    /tags/:name              
* Returns all posts with the given hashtag


