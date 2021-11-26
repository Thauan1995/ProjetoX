
$(document).ready(function(){
    $('#formulario-cadastro').on('submit', criarUsuario);
});

function criarUsuario(evento){
    evento.preventDefault();
    console.log("Dentro da função usuario");

    if ($('#senha').val() != $('#confirmar-senha').val()) {
        Swal.fire(
            'Viixe!',
            'As senhas não condizem!',
            'error'
        );
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

        Swal.fire(
            'Bem-vindo!',
            'Usuário cadastrado com sucesso!',
            'success'
        );

    }).fail(function(err) {
        console.log(err);
        Swal.fire(
            'Viixe!',
            'Error: ' + err,
            'error'
        );
    });
}