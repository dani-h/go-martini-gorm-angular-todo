/* global angular */
var app = angular.module("app", [])

app.controller("main", ["$scope", "$http", "$httpParamSerializerJQLike",
  function($scope, $http, $httpParamSerializerJQLike) {

    function addTodo(text, completed) {
      var $deferred = $http({
        url: "/todos/",
        method: "POST",
        data: $httpParamSerializerJQLike({
          text: text,
          completed: completed
        }),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      })

      $deferred.then(function(response) {
        $scope.todos.push(response.data)
        console.log($scope.todos.length)
      })
    }

    function setTodoCompletedState(todo) {
      $http({
        url: "/todos/" + todo.id,
        method: "PUT",
        data: $httpParamSerializerJQLike({
          text: todo.text,
          completed: todo.completed
        }),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      })
    }


    $scope.currentTodoInput = ""

    $scope.todoInputKeyup = function($event) {
      if ($event.keyCode === 13) {
        addTodo($scope.currentTodoInput, false)
      }
    }

    $scope.todoBtnClick = function() {
      if ($scope.currentTodoInput.length > 0) {
        addTodo($scope.currentTodoInput, false)
      }
    }

    $scope.todoCheckboxChange = function(todo) {
      setTodoCompletedState(todo)
    }

    $scope.clearCompleted = function() {
      $scope.todos.forEach(function(todo) {

      })
    }

    $http.get("/todos/")
      .then(function(response) {
        var todos = response.data
        $scope.todos = todos
      })
  }
])
