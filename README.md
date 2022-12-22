# GoNews
CIS 193 Final Project

---

## Description
This is a CRUD REST API for a blog-style website where users can create posts and add hashtags.
By adding #somehashtag to the post content, 'somehashtag' will be added to the post tags and a tag 
named 'somehashtag' will be created with an array of post IDs. 

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


