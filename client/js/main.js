function addUser(e) {

    var xhttp = new XMLHttpRequest();
    
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            getUsers();
        }
    };
    
    var name = document.getElementById("name").value;
    var email = document.getElementById("email").value;
    var password = document.getElementById("password").value;

    var data = { name: name, email: email, password: password };

    xhttp.open("POST", "/users/", true);
    xhttp.send(JSON.stringify(data));
}

function getUsers() {
    var xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {  
            document.getElementById("userTable").innerHTML = "";
            var users = JSON.parse(this.responseText);
            for(var i=0 ;i < users.length ; i++)
            {
                var a = users[i];
                var row = (document.createElement("tr"));
                var name = document.createElement("td");
                name.innerHTML = a.Name;
                var email = document.createElement("td");
                email.innerHTML = a.Email;
                row.appendChild(name);
                row.appendChild(email);
                document.getElementById("userTable").appendChild(row);
            }
            
            
            userCountElement = document.getElementById("user-count");
            userCountElement.innerHTML = users.length
        }
    };

    xhttp.open("GET", "/users/", true);
    xhttp.send();
}


; (function () {

    'use strict';

    var carousels = function () {
        jQuery(".owl-carousel1").owlCarousel(
            {
              loop:true,
              center: true,
              margin:0,
              responsiveClass:true,
              nav:false,
              responsive:{
                  0:{
                      items:1,
                      nav:false
                  },
                  600:{
                      items:1,
                      nav:false
                  },
                  1000:{
                      items:1,
                      nav:true,
                      loop:false
                  }
              }
          }
          );
        
          jQuery(".owl-carousel2").owlCarousel(
            {
              loop:true,
              center: true,
              margin:30,
              responsiveClass:true,
              nav:true,
              responsive:{
                  0:{
                      items:1,
                      nav:true
                  },
                  600:{
                      items:2,
                      nav:true,
                      margin:10,
                      center: false,
                  },
                  1000:{
                      items:3,
                      nav:true,
                      loop:true
                  }
              }
          }
          );
    }


    var isotope = function () {
        var $container = $('.portfolioContainer');
        $container.isotope({
            filter: '*',
            animationOptions: {
                duration: 750,
                easing: 'linear',
                queue: false
            }
        });

        $('.portfolioFilter a').click(function () {
            $('.portfolioFilter .active').removeClass('active');
            $(this).addClass('active');

            var selector = $(this).attr('data-filter');
            $container.isotope({
                filter: selector,
                animationOptions: {
                    duration: 750,
                    easing: 'linear',
                    queue: false
                }
            });
            return false;
        }); 
    };

    var navbar = function () {
        $(window).scroll(function () {
            $("nav.navbar").offset().top > -70 ? $(".navbar-fixed-top").addClass("top-nav-collapse") : $(".navbar-fixed-top").removeClass("top-nav-collapse")
        }),
        $(function () {
            $("a.page-scroll").bind("click", function (a) { var o = $(this); $("html, body").stop().animate({ scrollTop: $(o.attr("href")).offset().top - 58 }, 1e3), a.preventDefault()
        })
        });
    };

    (function ($) {
        carousels();
        isotope();
        navbar();
    })(jQuery);


}());
