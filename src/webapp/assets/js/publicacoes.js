$('#nova-publicacao').on('submit', criarPublicacao);

$(document).on('click', '.curtir-publicacao', curtirPublicacao);
$(document).on('click', '.descurtir-publicacao', descurtirPublicacao);

$('#atualizar-publicacao').on('click', atualizarPublicacao);

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
        alert("Erro ao curtir publicação");
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
        alert("Erro ao curtir publicação");
        console.log(publicacaoId);
    }).always(function() {
        elementoClicado.prop('disabled', false);
    });
}

function atualizarPublicacao(evento) {
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
        alert("Publicação editada com sucesso");
    }).fail(function(){
        alert("Erro ao editar a publicação");
    }).always(function(){
        $('#atualizar-publicacao').prop('disabled', false);
    })
}