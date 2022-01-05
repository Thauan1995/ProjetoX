$('#parar-de-seguir').on('click', paraDeSeguir);
$('#seguir').on('click', seguir);
$('#editar-usuario').on('submit', editar);
$('#atualizar-senha').on('submit', atualizarSenha);
$('#deletar-usuario').on('click', deletarUsuario);

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
            nome: $('#nome').val(),
            email: $('#email').val(),
            nick: $('#nick').val(),
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

function atualizarSenha(evento){
    evento.preventDefault();

    if ($('#nova-senha').val() != $('#confirmar-senha').val()) {
        Swal.fire("Ops...", "As senhas não coincidem", "warning");
        return;
    }

    $.ajax({
        url: "/web/atualizar-senha",
        method: "PUT",
        data: {
            atual: $('#senha-atual').val(),
            nova: $('#nova-senha').val()
        }
    }).done(function(){
        Swal.fire("Sucesso!", "Senha atualizada com sucesso!", "success")
            .then(function(){
                window.location = "/web/perfil";
            })
    }).fail(function(){
        Swal.fire("Ops...", "Erro ao atualizar a senha!", "error");
    });
}

function deletarUsuario() {
    Swal.fire({
        title: "Atenção",
        text: "Tem certeza que deseja apagar a sua conta? Essa é uma ação irreversivel",
        showCancelButton: true,
        cancelButtonText: "Cancelar",
        icon: "warning"
    }).then(function(confirmacao){
        if (confirmacao.value) {
            $.ajax({
                url: "/web/deletar-usuario",
                method: "DELETE"
            }).done(function(){
                Swal.fire("Sucesso!", "Sua conta foi excluido com sucesso!", "success")
                    .then(function(){
                        window.location = "/web/logout";
                    })
            }).fail(function(){
                Swal.fire("Ops...", "Ocorreu um erro ao excluir sua conta!", "error");
            });
        }
    })
}