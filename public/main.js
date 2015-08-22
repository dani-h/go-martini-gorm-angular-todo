/* global angular */
var app = angular.module("app", [])

app.run(["$http", function($http) {
  //Angular doesn't encode shit like JQuery does as it does json encoding by default. I don't know how to deal with this
  //in Go
  $http.defaults.headers.common["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.post["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.put["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.patch["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
}])

app.controller("main", ["$scope", "$http", "$httpParamSerializerJQLike", "$q",
  function($scope, $http, $httpParamSerializerJQLike, $q) {

    function NewTodo(text, completed) {
      return {
        text: text,
        completed: completed
      }
    }

    function addTodo(todo) {
      var $deferred = $http.post("/apiv0/todos/", $httpParamSerializerJQLike(todo))

      $deferred.then(function(response) {
        $scope.todos.push(response.data)
        console.log($scope.todos.length)
      })
    }

    function setTodoCompletedState(todo) {
      $http.put("/apiv0/todos/" + todo.id, $httpParamSerializerJQLike(todo))
    }


    $scope.currentTodoInput = ""

    $scope.todoInputKeyup = function($event) {
      if ($event.keyCode === 13) {
        var todo = NewTodo($scope.currentTodoInput, false)
        addTodo(todo)
        $scope.currentTodoInput = ""
      }
    }

    $scope.todoBtnClick = function() {
      if ($scope.currentTodoInput.length > 0) {
        var todo = NewTodo($scope.currentTodoInput, false)
        addTodo(todo)
      }
    }

    $scope.todoCheckboxChange = function(todo) {
      setTodoCompletedState(todo)
    }

    $scope.clearCompleted = function() {
      var $promises = $scope.todos
        .filter(function(todo) {
          return todo.completed
        })
        .map(function(todo) {
          return $http.delete("/apiv0/todos/" + todo.id)
        })

      $q.all($promises).then(function(responses) {
        var todosToRemove = responses.filter(function(response) {
          return response.status === 200
        })
        .map(function(response) {
          return response.data
        })

        $scope.todos = $scope.todos.filter(function(todo) {
          for(var remove of todosToRemove) {
            if(todo.id === remove.id) {
              return false
            }
          }
          return true
        })

      })
    }

    //Initialize all todos
    $http.get("/apiv0/todos/")
      .then(function(response) {
        var todos = response.data
        $scope.todos = todos
      })
  }
])
