$('#login').on('submit', fazerLogin);

function fazerLogin(evento) {
    evento.preventDefault();

    $.ajax({
        url: "/web/login",
        method: "POST",
        data: {
            email: $("#email").val(),
            senha: $("#senha").val(),
        }
    }).done(function(){
        window.location = "/web/home";
    }).fail(function(){
        alert("Usuário ou senha inválidos!");
    });
}