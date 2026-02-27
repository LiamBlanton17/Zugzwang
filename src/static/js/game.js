// Game.js handles the game engine for the frontend display
// id=chess-board is the main dom element as the root

// Some global constants
const WHITE = 0;
const BLACK = 1;
const PAWN = 0;
const KNIGHT = 1;
const BISHOP = 2;
const ROOK = 3;
const QUEEN = 4;
const KING = 5;
const EMPTY = -1;

class Game {
    constructor(gameId, elId) {
        // Validate input
        if (!gameId) {
            throw new Error('Must provide a game id to Game class');
        }
        if (!elId) {
            throw new Error('Must provide an element id to Game class');
        }

        // Validate DOM element
        const el = $(elId);
        if (el.length === 0) {
            throw new Error(`Element ID (${elId}) does not find a DOM element`)
        }

        // Create object and print to screen
        this.el = el;
        this.gameId = gameId;
        this.board = [
            new Piece('R'), new Piece('N'), new Piece('B'), new Piece('Q'), new Piece('K'), new Piece('B'), new Piece('N'), new Piece('R'),
            new Piece('P'), new Piece('P'), new Piece('P'), new Piece('P'), new Piece('P'), new Piece('P'), new Piece('P'), new Piece('P'),
            new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '),
            new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '),
            new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '),
            new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '), new Piece(' '),
            new Piece('p'), new Piece('p'), new Piece('p'), new Piece('p'), new Piece('p'), new Piece('p'), new Piece('p'), new Piece('p'),
            new Piece('r'), new Piece('n'), new Piece('b'), new Piece('q'), new Piece('k'), new Piece('b'), new Piece('n'), new Piece('r'),
        ];
        this._draw();
    }

    start() {

    }

    move() {

    }

    offerDraw() {

    }

    _writeAction() {

    }

    _readAction() {

    }

    _draw() {
        this.el.html(`<table class="chess-board-tbl">
            <tbody>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
                <tr>
                    <td></td><td></td><td></td><td></td>
                    <td></td><td></td><td></td><td></td>
                </tr>
            </tbody>
        </table>`);
    }

}

class Piece {
    constructor(type) {
        const map = {
            'K': {color: WHITE, piece: KING},
            'Q': {color: WHITE, piece: QUEEN},
            'R': {color: WHITE, piece: ROOK},
            'B': {color: WHITE, piece: BISHOP},
            'N': {color: WHITE, piece: KNIGHT},
            'P': {color: WHITE, piece: PAWN},
            'k': {color: BLACK, piece: KING},
            'q': {color: BLACK, piece: QUEEN},
            'r': {color: BLACK, piece: ROOK},
            'b': {color: BLACK, piece: BISHOP},
            'n': {color: BLACK, piece: KNIGHT},
            'p': {color: BLACK, piece: PAWN},
            ' ': {color: EMPTY, piece: EMPTY}
        };

        // Validate input
        if (!type) {
            throw new Error('Piece must be supplied with a type');
        }
        if (!(type in map)) {
            console.log(type);
            throw new Error('Piece type not supported');
        }

        const {color, piece} = map[type];
        this.color = color;
        this.piece = piece;
        this.type = type;
    }

    // Return HTML to place inside of a sqaure
    // Don't set color bg or anything, and eventually return an image
    draw() {
        // this.type is safe as it must be in the map for the object to be created
        return `<span class='piece'>${this.type}</span>`;
    }

}
