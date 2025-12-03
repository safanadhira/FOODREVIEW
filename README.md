Names :
1. Sanie Ghanda
2. Safa Nadhira

About Our Project : Food Review Website
Youtube Link : 

This project is a complete website where people can look at different restaurants, check out their menus, and write reviews for the food they eat. We built the whole system using the Go language and the Gin Framework for the back end, which is very fast and reliable. For storing all the information like restaurant names, menu items, and user reviews, we used a MySQL database managed by GORM. The website looks modern and easy to use because we designed the front end with a dark theme using Tailwind CSS.

The system is built on a clear structure, which is the Model-View-Controller (MVC) pattern. This helps keep the code organized and simple to manage. The key parts of the system are three main types of data that are all connected: a Restaurant can have many Food items (the menu), and each Food item can have many Reviews from users. To make sure the data stays clean, we added important rules in the database models. For example, if we delete a restaurant, the system automatically deletes all the related food items and their reviews. This is called Cascading Delete, and it prevents the database from getting messy with old, unlinked information.

A major feature of this application is how it handles user-uploaded pictures. We created special code to manage the entire process for images, like saving them safely to the server's storage and only keeping the picture's web address in the database. When a user updates a menu item with a new photo, the system is smart enough to delete the old picture file from the server to save space. Also, we made sure the website forms are secure and reliable. For instance, when someone writes a review, the system immediately checks that the rating is a number between 1 and 5 before saving it.

All parts of the website communicate smoothly through a map of web addresses called routes. This system follows a logical pattern (RESTful routing) to organize how the user moves through the site, whether they are looking at the main list of restaurants or going deep into a specific menu item to add a review. When displaying information, our backend is highly efficient: it loads a restaurant and all its menu items at once to speed up the page loading time. In short, this project is a robust, well-organized application that efficiently handles complex data relationships and user interactions.
