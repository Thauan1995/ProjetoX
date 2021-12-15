$('#nova-publicacao').on('submit', criarPublicacao);

$(document).on('click', '.curtir-publicacao', curtirPublicacao);
$(document).on('click', '.descurtir-publicacao', descurtirPublicacao);

$('#atualizar-publicacao').on('click', atualizarPublicacao);
$('.deletar-publicacao').on('click', deletarPublicacao);

function criarPublicacao(evento) {
    evento.preventDefault();

    $.ajax({
        url: "/web/publicacoes",
        method: "POST",
        data: {
            titulo: $('#titulo').val(),
            conteudo: $('#conteudo').val(),
        }
    }).done(function() {
        window.location = "/web/home";
    }).fail(function() {
        alert("Erro ao criar a publicação!");
    });
}

function curtirPublicacao(evento) {
    evento.preventDefault();

    const elementoClicado = $(evento.target);
    const publicacaoId = elementoClicado.closest('div').data('publicacao-id');

    elementoClicado.prop('disabled', true);
    $.ajax({
        url: `/web/publicacoes/${publicacaoId}/curtir`,
        method: "POST",
    }).done(function(){
        const contadorDeCurtidas = elementoClicado.next('span');
        const quantidadeDeCurtidas = parseInt(contadorDeCurtidas.text());

        contadorDeCurtidas.text(quantidadeDeCurtidas + 1);

        elementoClicado.addClass('descurtir-publicacao');
        elementoClicado.addClass('text-danger');
        elementoClicado.removeClass('curtir-publicacao');

    }).fail(function(){
        Swal.fire(
            'Ops...',
            'Erro ao criar publicação!',
            'error'
        );
        console.log(publicacaoId);
    }).always(function() {
        elementoClicado.prop('disabled', false);
    });
}

function descurtirPublicacao(evento) {
    evento.preventDefault();

    const elementoClicado = $(evento.target);
    const publicacaoId = elementoClicado.closest('div').data('publicacao-id');

    elementoClicado.prop('disabled', true);
    $.ajax({
        url: `/web/publicacoes/${publicacaoId}/descurtir`,
        method: "PUT",
    }).done(function(){
        const contadorDeCurtidas = elementoClicado.next('span');
        const quantidadeDeCurtidas = parseInt(contadorDeCurtidas.text());

        contadorDeCurtidas.text(quantidadeDeCurtidas - 1);

        elementoClicado.removeClass('descurtir-publicacao');
        elementoClicado.removeClass('text-danger');
        elementoClicado.addClass('curtir-publicacao');

    }).fail(function(){
        Swal.fire(
            'Ops...',
            'Erro ao curtir publicação!',
            'error'
        );
        console.log(publicacaoId);
    }).always(function() {
        elementoClicado.prop('disabled', false);
    });
}

function atualizarPublicacao() {
    $(this).prop('disebled', true);

    const publicacaoId = $(this).data('publicacao-id');
    
    $.ajax({
        url:`/web/publicacoes/${publicacaoId}`,
        method: "PUT",
        data: {
            titulo: $('#titulo').val(),
            conteudo: $('#conteudo').val()
        }
    }).done(function(){
        Swal.fire(
            'Sucesso!',
            'Publicação Atualizada com sucesso!',
            'success'
        )
        .then(function(){
            window.location = "/web/home";
        })
    }).fail(function(){
        Swal.fire(
            'Ops...',
            'Erro ao editar publicação!',
            'error'
        );;
    }).always(function(){
        $('#atualizar-publicacao').prop('disabled', false);
    })
}

function deletarPublicacao(evento) {
    evento.preventDefault();

    Swal.fire({
        title: "Atenção!",
        text: "Tem certeza que deseja excluir essa publicação? Essa ação é irreversível!",
        showCancelButton: true,
        cancelButtonText: "Cancelar",
        icon: "warning"
    }).then(function(confirmacao){
        if (!confirmacao.value) return;

        const elementoClicado = $(evento.target);
        const publicacao = elementoClicado.closest('div')
        const publicacaoId = publicacao.data('publicacao-id');
    
        elementoClicado.prop('disabled', true);
    
        $.ajax({
            url: `/web/publicacoes/${publicacaoId}/deletar`,
            method: "DELETE"
        }).done(function(){
            publicacao.fadeOut("slow", function(){
                $(this).remove();
            });
        }).fail(function(){
            Swal.fire(
                'Ops...',
                'Erro ao excluir publicação!',
                'error'
            );
        });
    })

}