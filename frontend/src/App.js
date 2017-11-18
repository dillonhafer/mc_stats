import React, { Component } from 'react';
import './App.css';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import statMap from './statMap';

const leaderKey = 220;
const DataRow = ({ label, value }) => {
  let v = value || '😐';

  if (label === 'Minutes Played') {
    v = (value * 0.05 / 60).toFixed(2);
    v += ' (' + (v / 60).toFixed() + ' hours)';
  }

  return (
    <tr>
      <td style={{ fontWeight: 'bold', color: '#777' }}>{label}</td>
      <td style={{ textAlign: 'right', fontWeight: 'bold' }}>{v}</td>
    </tr>
  );
};

const Player = ({ player, data }) => {
  return (
    <div className="player">
      <table style={{ width: '100%', padding: '8px' }}>
        <tbody>
          <DataRow label="Died:" value={data['stat.deaths']} />
          <DataRow label="Chests Opened:" value={data['stat.chestOpened']} />
          <DataRow label="Animals Bred:" value={data['stat.animalsBred']} />
          <DataRow label="Left Game:" value={data['stat.leaveGame']} />
          <DataRow
            label="Villagers Talked To:"
            value={data['stat.talkedToVillager']}
          />
          {Object.entries(statMap).map(e => {
            const label = e[1];
            const key = e[0];
            return <DataRow key={key} label={label} value={data[key]} />;
          })}
        </tbody>
      </table>
    </div>
  );
};

class App extends Component {
  state = {
    loading: false,
    players: [],
    stats: [],
    debug: false,
    requests: 0,
    lastUpdated: new Date(),
  };

  componentDidMount() {
    document.addEventListener('keydown', this.debugMode);
    this.loadPlayers();
    this.loadStats();
    window.setInterval(this.loadStats, 5000);
  }

  debugMode = e => {
    if (e.keyCode === leaderKey) {
      this.setState({ debug: !this.state.debug });
    }
  };

  renderDebug = () => {
    const { requests, lastUpdated } = this.state;
    return (
      <div className="debug">
        Debug mode
        <p>Requests Made: {requests}</p>
        <p>Last Updated: {lastUpdated.toString()}</p>
        <table>
          <tbody>
            <tr>
              <td colSpan="2">Players</td>
            </tr>
            {this.state.players.map(p => {
              return (
                <tr
                  key={p.uuid}
                  onClick={() => this.setState({ selectedDebugUUID: p.uuid })}
                >
                  <td>{p.name}</td>
                  <td>{p.uuid}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
        <button
          className="clearDebug"
          onClick={() => this.setState({ selectedDebugUUID: null })}
        >
          Clear Selected UUID
        </button>
        <hr />
        {this.state.selectedDebugUUID && this.renderDebugSelected()}
      </div>
    );
  };

  renderDebugSelected = () => {
    const { selectedDebugUUID, stats, debugSearch } = this.state;
    const data = stats.find(s => {
      return s.UUID === selectedDebugUUID;
    });

    return (
      <div>
        <input
          onChange={e => this.setState({ debugSearch: e.target.value })}
          type="search"
          placeholder="Search"
        />
        <table>
          <tbody>
            <tr>
              <td colSpan="2">Data for {data.UUID}</td>
            </tr>
            {Object.entries(data.data).map(e => {
              if (debugSearch) {
                if (e[0].match(new RegExp(`.*${debugSearch}.*`, 'gi'))) {
                  return (
                    <tr key={e[0]}>
                      <td>{e[0]}</td>
                      <td>{e[1]}</td>
                    </tr>
                  );
                } else {
                  return null;
                }
              }

              return (
                <tr key={e[0]}>
                  <td>{e[0]}</td>
                  <td>{e[1]}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  };

  request = async url => {
    try {
      this.setState({ loading: true });

      const request = {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        method: 'GET',
      };

      const r = await fetch(url, request);
      const json = await r.json();
      if (json) {
        return { ok: true, json };
      }
    } catch (err) {
      console.warn(err);
    } finally {
      this.setState({
        loading: false,
        requests: this.state.requests + 1,
        lastUpdated: new Date(),
      });
    }
  };

  loadStats = async () => {
    const resp = await this.request('/stats');
    if (resp && resp.ok) {
      this.setState({ stats: resp.json });
    }
  };

  loadPlayers = async () => {
    const resp = await this.request('/players');
    if (resp && resp.ok) {
      this.setState({ players: resp.json });
    }
  };

  players = () => {
    const { players, stats } = this.state;
    return players.map(p => {
      const stat = stats.find(s => {
        return s.UUID === p.uuid;
      }) || { data: {} };
      return (
        <TabPanel key={p.uuid}>
          <Player player={p} data={stat.data} />
        </TabPanel>
      );
    });
  };

  render() {
    const { loading, debug } = this.state;
    const notLoading = !loading;
    return (
      <div className="App">
        {debug && this.renderDebug()}
        {notLoading && this.players()}
        <Tabs defaultFocus={true}>
          <TabList>
            {this.state.players.map(({ name }) => <Tab key={name}>{name}</Tab>)}
          </TabList>
          {this.players()}
        </Tabs>
      </div>
    );
  }
}

export default App;
