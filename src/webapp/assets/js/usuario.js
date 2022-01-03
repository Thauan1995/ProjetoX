$('#parar-de-seguir').on('click', paraDeSeguir);
$('#seguir').on('click', seguir);
$('#editar-usuario').on('submit', editar);

function paraDeSeguir(){
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled',true);

    $.ajax({
        url: `/web/usuario/${usuarioId}/parar-de-seguir`,
        method: "POST"
    }).done(function(){
        window.location = `/web/usuario/${usuarioId}`;
    }).fail(function(){
        Swal.fire("Ops...","Erro ao parar de seguir o usuario!","error");
        $('#parar-de-seguir').prop('disabled', false);
    });
}

function seguir(){
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled',true);

    $.ajax({
        url: `/web/usuario/${usuarioId}/seguir`,
        method: "POST"
    }).done(function(){
        window.location = `/web/usuario/${usuarioId}`;
    }).fail(function(){
        Swal.fire("Ops...","Erro ao seguir o usuario!","error");
        $('#seguir').prop('disabled', false);
    });
}

function editar(evento){
    evento.preventDefault();

    $.ajax({
        url: "/web/editar-usuario",
        method: "PUT",
        data: {
            nome: $("#nome").val(),
            email: $("#email").val(),
            nick: $("#nick").val(),
        }
    }).done(function() {
        Swal.fire("Sucesso!", "Usuario atualizado com sucesso!", "success")
            .then(function() {
                window.location = "/web/perfil";
            });
    }).fail(function() {
        Swal.fire("Ops...", "Erro ao atualizar o usuário!", "error");
    });
}