/* global angular */
var app = angular.module("app", [])

app.run(["$http", ($http) => {
  //Angular doesn't encode shit like JQuery does as it does json encoding by default. I don't know how to deal with this
  //in Go
  $http.defaults.headers.common["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.post["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.put["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
  $http.defaults.headers.patch["Content-Type"] = 'application/x-www-form-urlencoded;charset=utf-8'
}])

app.controller("main", ["$scope", "$http", "$httpParamSerializerJQLike", "$q",
  function($scope, $http, $httpParamSerializerJQLike, $q) {

    const todosUrl = "/apiv0/todos/"

    function NewTodo(text, completed) {
      return {
        text: text,
        completed: completed
      }
    }

    function addTodo(todo) {
      var $deferred = $http.post(todosUrl, $httpParamSerializerJQLike(todo))
      $deferred.then((response) => $scope.todos.push(response.data))
    }

    function setTodoCompletedState(todo) {
      $http.put(todosUrl + todo.id, $httpParamSerializerJQLike(todo))
    }

    $scope.currentTodoInput = ""

    $scope.todoInputKeyup = ($event) => {
      if ($event.keyCode === 13) {
        var todo = NewTodo($scope.currentTodoInput, false)
        addTodo(todo)
        $scope.currentTodoInput = ""
      }
    }

    $scope.todoBtnClick = () => {
      if ($scope.currentTodoInput.length > 0) {
        var todo = NewTodo($scope.currentTodoInput, false)
        addTodo(todo)
      }
    }

    $scope.todoCheckboxChange = (todo) => {
      setTodoCompletedState(todo)
    }

    $scope.clearCompleted = () => {
      var $promises = $scope.todos
        .filter(todo => todo.completed)
        .map(todo => $http.delete(todosUrl + todo.id))

        $q.all($promises).then(responses => {

        var todosToRemove = responses
        .filter(response => response.status === 200)
        .map(response => response.data)

        $scope.todos = $scope.todos.filter((todo) => {
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
    $http.get(todosUrl)
      .then(response => $scope.todos = response.data)
  }
])
