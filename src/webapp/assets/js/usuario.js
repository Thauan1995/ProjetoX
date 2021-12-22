$('#parar-de-seguir').on('click', paraDeSeguir);
$('#seguir').on('click', seguir);

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