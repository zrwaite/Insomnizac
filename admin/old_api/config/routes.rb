Rails.application.routes.draw do
  resources :skills
  get "/users/login", to: "users#login"
  post "/users", to: "users#signup"
  post "/users/login", to: "users#login_handler"
  resources :users 

  resources :projects, param: :slug
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route ("/")
  root "home#index"
  get "/home", to: "home#index"
end