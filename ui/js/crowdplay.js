var crowdPlay = angular.module('crowdPlay', []);

crowdPlay.run(function($rootScope) {
    $rootScope.baseName= "CrowdPlay";
    $rootScope.channelSocket = io();
});

crowdPlay.controller('loginController', function($scope, $rootScope, $http, $window) {
    $scope.connectSocket = function(channelToJoin) {
        $rootScope.channelSocket.on('connect', function(response) {
           // Connected, let's sign-up for to receive messages for this room
           if (response) {
               console.log(response);
           }
           else {
               console.log("lel");
           }
           $rootScope.channelSocket.emit('joinChannel', channelToJoin);
        });
    };
});
