$('#formulario-cadastro').on('submit', criarUsuario);

function criarUsuario(evento){
    evento.preventDefault();
    console.log("Dentro da função usuario");

    if ($('#senha').val() != $('#confirmar-senha').val()) {
        alert("As senhas não coincidem!");
        return;
    }

    $.ajax({
        url: "/web/usuario/registrar",
        method: "POST",
        data: {
            nome: $('#nome').val(),
            nick: $('#nick').val(),
            email: $('#email').val(),
            senha: $('#senha').val()
        }
    }).done(function() {
        alert("Usuário cadastrado com sucesso!");
        window.location.href = "http://localhost:8000/web/login";
    }).fail(function(err) {
        console.log(err);
        alert("Erro ao cadastrar o usuário");
    });
}