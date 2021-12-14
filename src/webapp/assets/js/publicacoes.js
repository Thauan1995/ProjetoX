$('#nova-publicacao').on('submit', criarPublicacao);
$('.curtir-publicacao').on('click', curtirPublicacao);

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
    }).fail(function(){
        alert("Erro ao curtir publicação");
        console.log(publicacaoId);
    }).always(function() {
        elementoClicado.prop('disabled', false);
    })
}